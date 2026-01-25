import { useEffect, useState } from "react";
import "./App.css";  

type users = {
  id: number;
  username: string;
  created_at: string;
}

type topics = {
  id: number;
  title: string;
  description: string;
  user_id: number;
  created_at: string;
  updated_at: string | null;
}

type posts = {
  id: number;
  user_id: number;
  title: string;
  content: string;
  created_at: string;
}
type comments = {
  id: number;
  post_id: number;
  user_id: number;
  comment: string;
  created_at: string;
}

function MyButton() {
  return (
    <button>I'm a wizard</button>
  );
}

export default function App() {
  const [users, setUsers] = useState<users[]>([]);
  const [topics, setTopics] = useState<topics[]>([]);
  const [posts, setPosts] = useState<posts[]>([]);
  const [comments, setComments] = useState<comments[]>([]);

  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const [newTitle, setNewTitle] = useState("");
  const [newContent, setNewContent] = useState("");
  const [posting, setPosting] = useState(false);



  useEffect(() => {
   async function loadAll() {
    try {
      setLoading(true);
      setError(null);

      const [usersResponse, topicsResponse, postsResponse, commentsResponse] = await Promise.all([
        fetch("http://localhost:8080/users"),
        fetch("http://localhost:8080/topics"),
        fetch("http://localhost:8080/posts"),
        fetch("http://localhost:8080/comments"),
      ]);

      if (!usersResponse.ok || !topicsResponse.ok || !postsResponse.ok || !commentsResponse.ok) {
        throw new Error("Network response was not ok");
      }

      const [usersData, topicsData, postsData, commentsData] = await Promise.all([
        usersResponse.json(),
        topicsResponse.json(),
        postsResponse.json(),
        commentsResponse.json(),
      ]);

      setUsers(usersData);
      setTopics(topicsData);
      setPosts(postsData);
      setComments(commentsData);
    } catch (error) {
      console.error("Error fetching data:", error);
      setError("Failed to fetch data");
    } finally {
      setLoading(false);
    }
  }

  loadAll();

}, []);

//CREATE POST
async function createPost() {
  try {
    setPosting(true);

    const res = await fetch("http://localhost:8080/posts", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        user_id: 1, //to fix
        title: newTitle,
        content: newContent,
      }),
    });

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`);
    }

    setNewTitle("");
    setNewContent("");

    const updated = await fetch("http://localhost:8080/posts");
    setPosts(await updated.json());
  } catch (e) {
    console.error(e);
    alert("Failed to create post");
  } finally {
    setPosting(false);
  }
}

  return (
    <div>
      <h1>The Wizarding side of the Internet!!</h1>
      <h2>Welcome to all things magical</h2>
      <MyButton/>
      <h1>Explore the Wizarding World</h1>
      {loading && <p>Loading...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}

      {!loading && !error && (
        <div>
          <h2>Users</h2>
          <ul>
            {users.map((user) => (
              <li key={user.id}>{user.username}</li>
            ))}
          </ul>

          <h2>Topics</h2>
          <ul>
            {topics.map((topic) => (
              <li key={topic.id}>{topic.title}</li>
            ))}
          </ul>

          <h2>Posts</h2>
          <ul>
            {posts.map((post) => (
              <li key={post.id}>{post.title}</li>
            ))}
          </ul>

          <h2>Comments</h2>
          <ul>
            {comments.map((comment) => (
              <li key={comment.id}>{comment.comment}</li>
            ))}
          </ul>

          <section className="card">
  <h2>Create New Post</h2>

  <input
    placeholder="Post title"
    value={newTitle}
    onChange={(e) => setNewTitle(e.target.value)}
    style={{ width: "100%", marginBottom: 8 }}
  />

  <textarea
    placeholder="Post content"
    value={newContent}
    onChange={(e) => setNewContent(e.target.value)}
    style={{ width: "100%", marginBottom: 8 }}
  />

  <button className="btn" onClick={createPost} disabled={posting}>
    {posting ? "Posting..." : "Create Post"}
  </button>
</section>
        </div>
      )}
    </div>


  )
}
