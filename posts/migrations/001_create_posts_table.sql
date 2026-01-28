CREATE TABLE posts {
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category_id INT NOT NULL
    created_at TEXT DEFAULT (CURRENT_TIMESTAMP)
};

-- ajouter foreign key category_id