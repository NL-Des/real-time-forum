/* C'est une sécurité face à la suppression sauvage
et face à des comportements illogiques.  */
PRAGMA foreign_keys = ON;

/* Table des utilisateurs */
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    UserName TEXT NOT NULL UNIQUE CHECK(LENGTH(UserName) >= 3),
    Age INTEGER NOT NULL CHECK(Age >= 15), /* limite d'âge */
    Gender TEXT NOT NULL CHECK(Gender IN('Man', 'Woman', 'Other')), /* N'accepte que les termes indiqués */
    FirstName TEXT NOT NULL,
    LastName TEXT NOT NULL,
    Email TEXT NOT NULL UNIQUE,
    Password TEXT NOT NULL CHECK(LENGTH(Password) >= 8) /* Taille minimale du mdp */
);

/* Contient des comments */
CREATE TABLE IF NOT EXISTS post (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Title TEXT NOT NULL,
    Content TEXT NOT NULL,
    AuthorID INTEGER NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME,
    FOREIGN KEY(AuthorID) REFERENCES users(id)
);

/* Se trouve dans un post */
CREATE TABLE IF NOT EXISTS comments (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,    
    PostID INTEGER NOT NULL,
    AuthorID INTEGER NOT NULL,
    Content TEXT NOT NULL,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME,
    FOREIGN KEY (PostID) REFERENCES post (ID),
    FOREIGN KEY (AuthorID) REFERENCES users (id)
);

/* Table des messages du websocket */
CREATE TABLE IF NOT EXISTS messages (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    SenderID INTEGER NOT NULL,
    ReceiverID INTEGER NOT NULL,
    Content TEXT NOT NULL,
    WriteAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    ReadAt DATETIME,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt DATETIME,
    FOREIGN KEY (SenderID) REFERENCES users(id),
    FOREIGN KEY (ReceiverID) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS category (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
);

/* Table de liaison entre les tables "category" et "post" */
CREATE TABLE IF NOT EXISTS post_categories (
    PostID INTEGER NOT NULL, /* Colonne 1 */
    CategoryID INTEGER NOT NULL, /*  Colonne 2 */
    PRIMARY KEY (PostID, CategoryID), /* Clé primaire composée des colonnes précédentes */
	/* Grâce à la clé primaire, on ne peut pas avoir plusieurs fois la même catégorie pour un même post. */
    FOREIGN KEY (PostID) REFERENCES post(ID),
    FOREIGN KEY (CategoryID) REFERENCES category(ID)
);

CREATE TABLE IF NOT EXISTS session (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    UserID INTEGER NOT NULL,
    Token TEXT NOT NULL UNIQUE,
    CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP, /* Moment de la connexion */
    ExpiresAt DATETIME NOT NULL, /* 24 heures ? */
    UserAgent TEXT,
    IP TEXT,
    FOREIGN KEY (UserID) REFERENCES users(id)
);