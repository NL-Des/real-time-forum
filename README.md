# Real-Time-Forum

# Instructions
Objectif du projet, créer un forum en JS qui s'actualise en temps réel.

### Fonctionnalités Attendues :
- Enregistrement.
- Login.
- Création de posts.
- Commenter un post.
- Messages privés.

### Langages autorisés :

- SQLite pour stocker les données.
- Goland pour manipuler les données et les Websockets.(Backend)
    - Tous les packages standards (https://pkg.go.dev/std).
    - Websocket (https://pkg.go.dev/github.com/gorilla/websocket)
    - Sécurité chiffrement (https://pkg.go.dev/golang.org/x/crypto/bcrypt)
    - Sécurité connexion (https://github.com/gofrs/uuid ou https://github.com/google/uuid)
- Javascript pour le Frontend et les clients Websockets.
- HTML
- CSS

### Contraintes :
Une seule page HTML autorisée.

## Enregistrement et connexion
Les utiliseurs non connectés ne doivent que voir la page connexion/inscription.

### Informations nécessaires :
    - Pseudo
    - Âge
    - Genre
    - Nom
    - Prénom
    - Mail
    - Mot de passe

### Contraintes :
-L'utilisateur peut se connecter en utilisant son pseudo ou son mail avec le mot de passe?
-L'utilisateur doit pouvoir se déconnecter depuis n'importe quelle page du forum.

## Posts et commentaires :
L'utilisateur doit pouvoir :
- Créer des posts.
    - Mettre des catégories aux posts.
    - Créer des commentaires dans les posts.
- Voir les posts dans un fil d'actualité.
    - Voir les commentaires des posts si ils cliquent dessus.

## Messages privés :
Les utilisateurs doivent pouvoir s'envoyer des messages privés entre eux au travers d'un chat.
Deux sections sont attendues.

### Messages :
- Les messages doivent indiquer la date d'envoi.
- Les messages doivent indiquer le pseudo de l'envoyeur.
- L'utilisateur doit recevoir une notification

### Section liste des utilisateurs :
- Il doit être indiqué si l'utilisateur est connecté ou non.
- Les utilisateurs avec qui il y a eu les derniers messages échangés doivent êtres présentés en premiers.
    - Si l'utilisateur est nouveau et qu'il n'a pas encore envoyé/reçu de message, la liste des utilisateurs avec qui il peut échanger sera organisée par ordre alphabétique.
- L'utilisateur doit pouvoir envoyer des messages aux utilisateurs conenctés.
- La section doit toujours être visible par l'utilisateur.

### Section zone de messages avec les utilisateurs :
- Quand l'utilisateur clique sur un utilisateur pour lui envoyer un message, il doit charger les 10 derniers messages échangés.
    - Quand l'utilisateur scroll pour voir plus de messages, il doit en charger 10 de plus à chaque fois qu'il arrive au 10ème.
- L'utilisateur doit recevoir une notification à chaque fois qu'il reçoit un nouveau message. Sans à avoir à rafraichir la page.