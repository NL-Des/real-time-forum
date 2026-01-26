/* C'est une sécurité face à la suppression sauvage
et face à des comportements illogiques.  */
PRAGMA foreign_keys = ON;

/* Table des utilisateurs */
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    UserName TEXT NOT NULL UNIQUE CHECK(LENGTH(UserName) >= 3),
    Age INTEGER NOT NULL CHECK(Age => 15), /* Limite d'âge minimale */
	Gender TEXT NOT NULL CHECK(Gender IN('M', 'F', 'Other')), /* N'accepte que les informations entre '' */
	FirstName TEXT NOT NULL,
	LastName TEXT NOT NULL,
	Email TEXT NOT NULL UNIQUE,
	Password TEXT NOT NULL CHECK(LENGTH(Password) >= 8),
    userOnline INTEGER DEFAULT 0 /* 0 = hors ligne | 1 = en ligne */
);

/* Contient des comments */
CREATE TABLE post (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	Title TEXT NOT NULL,
	Content TEXT NOT NULL,
	AuthorID TEXT NOT NULL,
	CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	CategoryIDs INTEGER,
    FOREIGN KEY(AuthorID) REFERENCES users(UserName),
    FOREIGN KEY(CategoryIDs) REFERENCES category(ID)
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
	ReceiverID TEXT NOT NULL, /* Est-ce la bonne manière ? */
	Content TEXT NOT NULL,
	CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (SenderID) REFERENCES users(UserName)
	FOREIGN KEY (ReceiverID) REFERENCES users(UserName) /* Est-ce la bonne manière ? */
);

CREATE TABLE category (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL UNIQUE
)

/* Table de liaison entre les tables "category" et "post" */
CREATE TABLE post_categories (
	PostID INTEGER NOT NULL, /* définition de la colonne 1 */
	CategoryID INTEGER NOT NULL, /* définition de la colonne 2 */
	PRIMARY KEY (PostID, CategoryID) /* Clé primaire composite composée des colonnes précédentes*/
	/* Permet d'éviter d'avoir deux fois la même catégorie pour un même post. */
	FOREIGN KEY (PostID) REFERENCES post(id)
	FOREIGN KEY (CategoryID) REFERENCES category(ID)
)