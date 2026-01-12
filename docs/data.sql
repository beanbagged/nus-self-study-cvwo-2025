    INSERT INTO users (username) VALUES ('Harry'), ('Hermione'), ('Ron');

    INSERT INTO topics (title, description, user_id) VALUES 
    ('Magic', 'Discussion about magic spells and potions.', 1),
    ('Hogwarts', 'All about life at Hogwarts School of Witchcraft and Wizardry.', 2);

    INSERT INTO posts (user_id, title, content) VALUES 
    (1, 'Best Spells for Beginners', 'Let''s discuss the best spells for beginners in magic.'),
    (2, 'Hogwarts Houses', 'Which Hogwarts house do you belong to and why?');

    INSERT INTO comments (post_id, user_id, comment) VALUES 
    (1, 2, 'I think Wingardium Leviosa is a great spell to start with!'),
    (2, 3, 'I belong to Gryffindor because of my bravery.');

    SELECT * FROM posts WHERE user_id = (SELECT id FROM users WHERE username = 'Hermione');
