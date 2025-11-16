# Al - CLI de gestion de projets

> **Al** est un outil en ligne de commande pour gérer efficacement vos projets depuis le terminal.

## 🎯 À quoi ça sert ?

Vous jonglez avec plusieurs projets clients, chacun dans un dossier différent ? Vous devez :
- Naviguer rapidement entre vos projets
- Stocker des notes techniques ou sensibles (mots de passe, clés API)
- Conserver des liens utiles (documentation, repos, dashboards)
- Avoir des raccourcis pour retrouver vos projets instantanément

**Al** centralise tout ça en ligne de commande pour un accès ultra-rapide.

## 🗂️ Les objets manipulés

### 1. **Projets**
Un **projet** est un répertoire de travail identifié par :
- Son **nom** (le nom du dossier)
- Des **shortcuts** (raccourcis personnalisés pour y accéder rapidement)
- Son **chemin absolu** sur le système

**Exemple** : Le projet dans `/home/alex/client-acme` avec les shortcuts `acme`, `client1`, `ca`

### 2. **Notes**
Les **notes** sont des morceaux de texte attachés à un projet :
- Stockées localement dans `.al_local/notes/`
- Peuvent être **chiffrées** (AES-GCM) pour les informations sensibles
- Éditables avec vim ou directement en ligne de commande
- Copiables dans le presse-papier

**Cas d'usage** : Notes de réunion, credentials, commandes fréquentes, TODO techniques

### 3. **Links**
Les **links** sont des URLs avec des mots-clés :
- Nom principal et keywords multiples pour retrouver facilement
- Stockés dans `.al_local/links/`
- Copiables dans le presse-papier

**Cas d'usage** : Documentation, repos GitHub, dashboards de monitoring, outils en ligne

## 🏗️ Architecture

```
~/.al_global/              # Configuration globale
├── projects               # Registry de tous les projets
└── config                 # Paramètres (longueur préview, etc.)

<projet>/.al_local/        # Données locales du projet
├── notes/                 # Notes (JSON avec contenu chiffré ou non)
└── links/                 # Links (JSON avec URL et keywords)
```

## 📦 Installation

### 1. Cloner le repository

```bash
git clone https://github.com/agourdel/al.git
cd al
```

### 2. Télécharger les dépendances

```bash
make deps
```

### 3. Compiler les binaires

```bash
make build
```

### 4. Installer sur le système

```bash
sudo make install
```

Cela installe les binaires suivants dans `/usr/local/bin` :
- `al` - Commande principale
- `algo` - Raccourci pour `al go`
- `alinit` - Raccourci pour `al init`
- `alnote` - Raccourci pour `al note`
- `allink` - Raccourci pour `al link`

### 5. Mettre à jour (après modifications)

```bash
make build
sudo build/al update
```

## 📖 Commandes disponibles

### 🔧 Commandes système

#### `al install`
Installe la CLI sur le système (à faire une seule fois).

```bash
al install
```

#### `al update`
Met à jour les binaires après une recompilation.

```bash
al update
```

---

### 📁 Gestion des projets

#### `al init [shortcuts...]`
Initialise le répertoire courant comme projet Al.

```bash
# Avec le nom du dossier uniquement
al init

# Avec des raccourcis supplémentaires (séparés par espaces)
al init foo bar test

# Ou avec des pipes (entre guillemets)
al init "client1|c1|acme"
alinit "ways|wayz|wa"
```

**Résultat** :
- Crée `.al_local/` dans le répertoire courant
- Enregistre le projet dans `~/.al_global/projects`
- Tous les shortcuts sont utilisables pour retrouver le projet

#### `al go <shortcut>` ou `algo <shortcut>`
Copie le chemin absolu du projet dans le presse-papier.

```bash
algo myproject    # Copie le chemin
cd <Ctrl+V>       # Colle et navigue
```

**Astuce** : Si le shortcut n'existe pas, la CLI suggère des noms similaires.

---

### 📝 Gestion des notes

#### `al note list` ou `alnote list`
Liste toutes les notes du projet courant.

```bash
alnote list
alnote list -t autre_projet    # Pour un autre projet
```

Affiche un tableau :
```
Name        Date        Preview
----        ----        -------
todo        2025-11-16  Faire la review du code...
api_key     2025-11-15  **chiffrée**
```

#### `al note add #nom` ou `alnote add #nom`
Ajoute une nouvelle note.

```bash
# Note simple (ouvre vim)
alnote add #todo

# Note chiffrée (demande un mot de passe)
alnote add #api_key -c

# Note avec contenu direct (sans éditeur)
alnote add #reminder -b "Ne pas oublier de push"

# Note chiffrée pour un autre projet
alnote add #secret -t client2 -c -b "Password: 123456"
```

**Options** :
- `-c, --chiffre` : Chiffre la note (AES-GCM avec mot de passe)
- `-b, --body <text>` : Contenu direct sans ouvrir l'éditeur
- `-t, --target <project>` : Cibler un autre projet

#### `al note get #nom` ou `alnote get #nom`
Affiche ou copie une note.

```bash
# Afficher dans le terminal
alnote get #todo

# Copier dans le presse-papier
alnote get #api_key --cp

# Note chiffrée (demande le mot de passe)
alnote get #secret -c
```

#### `al note edit #nom` ou `alnote edit #nom`
Modifie une note existante.

```bash
# Ouvrir dans vim
alnote edit #todo

# Remplacer directement le contenu
alnote edit #reminder -b "Nouveau contenu"
```

#### `al note remove #nom` ou `alnote remove #nom`
Supprime une note (avec confirmation).

```bash
alnote remove #old_note
alnote remove #secret -t autre_projet
```

---

### 🔗 Gestion des liens

#### `al link list` ou `allink list`
Liste tous les liens du projet.

```bash
allink list
allink list -t autre_projet
```

Affiche un tableau :
```
Name      Link                              Keywords
----      ----                              --------
github    https://github.com/user/repo      repo, code, git
docs      https://docs.example.com          documentation, aide
```

#### `al link add #nom -u <url>` ou `allink add #nom -u <url>`
Ajoute un nouveau lien.

```bash
allink add #github -u https://github.com/user/repo
allink add #docs -u https://docs.example.com -k "documentation|aide|help"
allink add #api -u https://api.example.com -k "api|rest"
```

**Options** :
- `-u, --url <url>` : URL du lien (obligatoire)
- `-k, --keywords <keywords>` : Mots-clés séparés par `|`
- `-t, --target <project>` : Cibler un autre projet

#### `al link get #nom` ou `allink get #nom`
Affiche ou copie un lien (recherche par nom ou keyword).

```bash
# Afficher les détails
allink get #github
allink get documentation    # Recherche par keyword

# Copier l'URL dans le presse-papier
allink get #api --cp
allink get rest --cp        # Par keyword
```

#### `al link edit #nom` ou `allink edit #nom`
Modifie un lien existant.

```bash
# Changer l'URL
allink edit #github -u https://github.com/newuser/repo

# Ajouter des keywords
allink edit #docs -ak "tutorial|tuto"

# Réinitialiser les keywords
allink edit #api -rk "rest|api|v2"

# Supprimer tous les keywords
allink edit #api -rk ""
```

**Options** :
- `-u, --url <url>` : Nouvelle URL
- `-ak, --add_keyword <keywords>` : Ajouter des keywords
- `-rk, --reset_keyword <keywords>` : Remplacer tous les keywords

#### `al link remove #nom` ou `allink remove #nom`
Supprime un lien (avec confirmation).

```bash
allink remove #github
allink remove api    # Par keyword
```

---

## 🎯 Fonctionnalités clés

### 🔐 Chiffrement des notes
Les notes sensibles sont chiffrées avec **AES-GCM** :
- Dérivation de clé avec PBKDF2 (4096 itérations)
- Chaîne de vérification intégrée pour détecter un mauvais mot de passe
- Sel aléatoire unique par note

### 🔍 Suggestions intelligentes
Quand un nom n'est pas trouvé, la CLI calcule la **distance de Levenshtein** et suggère des noms similaires :

```bash
$ algo wayys
Project 'wayys' not found. Did you mean:
  - waays
  - wayz
  - ways
```

### 📋 Presse-papier
Intégration native avec le presse-papier système :
- Notes : `--cp` pour copier le contenu
- Links : `--cp` pour copier l'URL
- Projects : `algo` copie automatiquement le chemin

### 🎨 Projet distant avec `-t`
Toutes les commandes `note` et `link` supportent `-t` pour cibler un autre projet :

```bash
alnote list -t client2
allink add #dashboard -u https://dash.client2.com -t client2
```

---

## 🚀 Commandes futures possibles

### Gestion avancée des projets
- `al list` : Lister tous les projets enregistrés avec leurs shortcuts
- `al rename <old> <new>` : Renommer un projet
- `al archive <project>` : Archiver un projet (le garder en registry mais le marquer comme archivé)
- `al unarchive <project>` : Réactiver un projet archivé
- `al remove <project>` : Supprimer complètement un projet du registry
- `al info <project>` : Afficher toutes les infos d'un projet (path, shortcuts, nombre de notes/links)
- `al sync` : Synchroniser les projets (vérifier que les chemins existent toujours)

### Commandes rapides
- `al cmd add #name <command>` : Sauvegarder des commandes shell fréquentes
- `al cmd run #name` : Exécuter une commande sauvegardée
- `al cmd list` : Lister les commandes
- `alcmd` : Alias pour `al cmd`

### Variables d'environnement
- `al env add <KEY=VALUE>` : Sauvegarder des variables d'env par projet
- `al env list` : Lister les variables
- `al env export` : Exporter toutes les variables dans le shell actuel
- `al env load` : Charger les variables (génère un script sourceable)

### Tags et filtres
- `al tag <project> <tag>` : Ajouter un tag à un projet
- `al list --tag <tag>` : Filtrer les projets par tag
- `alnote add #note --tags important,urgent` : Ajouter des tags aux notes
- `alnote list --tag urgent` : Filtrer les notes par tag

### Import/Export
- `al export <project>` : Exporter un projet (notes + links) en JSON/YAML
- `al import <file>` : Importer un projet depuis un fichier
- `al backup` : Backup de tous les projets
- `al restore <backup>` : Restaurer depuis un backup

### Recherche globale
- `al search <query>` : Rechercher dans toutes les notes et links de tous les projets
- `al search <query> -t <project>` : Rechercher dans un projet spécifique
- `al search <query> --encrypted` : Inclure les notes chiffrées (demande les mots de passe)

### Collaboration
- `al share <project>` : Générer un fichier de partage (sans notes chiffrées)
- `al team add <email>` : Ajouter un membre à un projet (nécessiterait un backend)
- `al remote add <url>` : Synchroniser avec un serveur distant

### Templates
- `al template create <name>` : Créer un template de projet avec notes/links prédéfinis
- `al template use <name>` : Initialiser un projet depuis un template
- `al template list` : Lister les templates disponibles

### Statistiques
- `al stats` : Afficher des stats (nombre de projets, notes, links)
- `al stats <project>` : Stats détaillées d'un projet
- `al timeline` : Timeline des modifications récentes

### Intégrations
- `al github link <repo>` : Lier automatiquement un repo GitHub
- `al jira link <project>` : Lier un projet Jira
- `al slack notify <message>` : Envoyer une notification Slack
- `al notion export` : Exporter vers Notion

---

## 💻 Développement

### Compiler
```bash
make build
```

### Nettoyer
```bash
make clean
```

### Tests (à implémenter)
```bash
make test
```

---

## 📄 License

MIT

---

## 🤝 Contribution

Les contributions sont les bienvenues ! N'hésitez pas à ouvrir une issue ou une pull request.
