PRAGMA foreign_keys = ON;

/* Table des utilisateurs */
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    UserName TEXT NOT NULL UNIQUE CHECK(LENGTH(UserName) >= 3),
    Age INTEGER NOT NULL CHECK(Age >= 15), /* CORRIGÉ: >= au lieu de => */
    Gender TEXT NOT NULL CHECK(Gender IN('Man', 'Woman', 'Other')),
    FirstName TEXT NOT NULL,
    LastName TEXT NOT NULL,
    Email TEXT NOT NULL UNIQUE,
    Password TEXT NOT NULL CHECK(LENGTH(Password) >= 8),
    userOnline INTEGER DEFAULT 0
);

/* Contient des comments */
CREATE TABLE post (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Title TEXT NOT NULL,
    Content TEXT NOT NULL,
    AuthorID TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    /* CategoryIDs supprimé car vous utilisez post_categories */
    FOREIGN KEY(AuthorID) REFERENCES users(UserName)
);

/* Se trouve dans un post */
CREATE TABLE comments (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,    
    PostID INTEGER NOT NULL,
    AuthorID TEXT NOT NULL,
    Content TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (PostID) REFERENCES post (ID),
    FOREIGN KEY (AuthorID) REFERENCES users (UserName)
);

/* Table des messages du websocket */
CREATE TABLE messages (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    SenderID TEXT NOT NULL,
    ReceiverID TEXT NOT NULL,
    Content TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (SenderID) REFERENCES users(UserName),
    FOREIGN KEY (ReceiverID) REFERENCES users(UserName)
);

CREATE TABLE category (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
);

/* Table de liaison entre les tables "category" et "post" */
CREATE TABLE post_categories (
    PostID INTEGER NOT NULL,
    CategoryID INTEGER NOT NULL,
    PRIMARY KEY (PostID, CategoryID),
    FOREIGN KEY (PostID) REFERENCES post(ID),
    FOREIGN KEY (CategoryID) REFERENCES category(ID)
);