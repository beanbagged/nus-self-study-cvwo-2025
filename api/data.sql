    INSERT INTO users (username) VALUES ('Harry'), ('Hermione'), ('Ron');

    INSERT INTO topics (title, description, user_id) VALUES 
    ('Magic', 'Discussion about spells and potions.', 1),
    ('Hogwarts', 'All about campus life at Hogwarts!', 2);

    INSERT INTO posts (user_id, title, content) VALUES 
    (1, 'Best Spells for Beginners', 'Let''s discuss the best spells to start with.'),
    (2, 'Hogwarts Houses', 'Which is the best Hogwarts house?');

    INSERT INTO comments (post_id, user_id, comment) VALUES 
    (1, 2, 'I think Wingardium Leviosa is easy enough to pick up.'),
    (2, 3, 'Gryffindor because we are brave!');

    SELECT * FROM posts WHERE user_id = (SELECT id FROM users WHERE username = 'Hermione');