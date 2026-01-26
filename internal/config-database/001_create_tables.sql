/* C'est une sécurité face à la suppression sauvage
et face à des comportements illogiques.  */
PRAGMA foreign_keys = ON;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    UserName TEXT NOT NULL UNIQUE,
    Age INTEGER NOT NULL,
	Gender TEXT NOT NULL,
	FirstName TEXT NOT NULL,
	LastName TEXT NOT NULL,
	Email TEXT NOT NULL,
	Password TEXT NOT NULL,
    userOnline INTEGER DEFAULT 0 /* 0 = hors ligne | 1 = en ligne */
    /* Question importante */
    /* Si le user est supprimé, il y a un delete en cascade ? */
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
    FOREIGN KEY(AuthorID) REFERENCES users(id),
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
    FOREIGN KEY (AuthorID) REFERENCES users (id)
);

CREATE TABLE messages (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	SenderID TEXT NOT NULL,
	ReceiverID /* Comment l'obtenir, A lier */
	Content TEXT NOT NULL,
	CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
	UpdatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (SenderID) REFERENCES users(id)
);

CREATE TABLE category (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Name TEXT NOT NULL
    /* Si la catégorie est supprimée, il y a un delete en cascade. */
)