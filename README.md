# Al - CLI de gestion de projets clients

Al est un outil CLI en Golang pour gérer vos projets clients avec des notes, des liens et des raccourcis.

## Installation

1. Cloner le repository et naviguer dans le dossier
2. Télécharger les dépendances :
   ```bash
   make deps
   ```
3. Compiler les binaires :
   ```bash
   make build
   ```
4. Installer sur le système (nécessite sudo) :
   ```bash
   make install
   ```

Cela installera les binaires suivants dans `/usr/local/bin` :
- `al` - Commande principale
- `algo` - Raccourci pour `al go`
- `alinit` - Raccourci pour `al init`
- `alnote` - Raccourci pour `al note`
- `allink` - Raccourci pour `al link`



## Commandes

### `al install`
Installe la CLI sur le système. Crée le répertoire global `~/.al_global` et copie les binaires dans `/usr/local/bin`.

```bash
al install
```

### `al update`
Met à jour les binaires existants sur le système.

```bash
al update
```

### `al init [shortcuts...]`
Initialise le répertoire courant comme projet Al. Le nom du dossier est automatiquement utilisé comme raccourci, et vous pouvez en ajouter d'autres.

```bash
# Initialiser avec le nom du dossier seulement
al init

# Initialiser avec des raccourcis supplémentaires
al init foo|bar|test
alinit client1|c1
```

### `al go <shortcut>`
Navigue vers le répertoire d'un projet en utilisant son nom ou un raccourci.

```bash
al go myproject
algo myproject
algo foo
```

### `al note`
Gère les notes pour un projet.

#### Lister les notes
```bash
al note list
al note list -t autre_projet
alnote list
```

#### Ajouter une note
```bash
# Note simple avec éditeur
al note add #ma_note

# Note chiffrée
al note add #secret -c

# Note avec contenu direct
al note add #todo -b "Faire ceci et cela"

# Note chiffrée pour un autre projet
al note add #secret -t projet2 -c -b "Contenu secret"
```

#### Obtenir une note
```bash
# Afficher dans le terminal
al note get #ma_note

# Copier dans le presse-papier
al note get #ma_note --cp

# Avec un autre projet
alnote get #secret -t projet2
```

#### Éditer une note
```bash
# Ouvrir dans l'éditeur
al note edit #ma_note

# Remplacer le contenu
al note edit #ma_note -b "Nouveau contenu"
```

#### Supprimer une note
```bash
al note remove #ma_note
alnote remove #secret -t autre_projet
```

### `al link`
Gère les liens HTTPS pour un projet.

#### Lister les liens
```bash
al link list
al link list -t autre_projet
allink list
```

#### Ajouter un lien
```bash
al link add #github -u https://github.com/user/repo
al link add #docs -u https://docs.example.com -k documentation|aide|help
allink add #api -u https://api.example.com -k api|rest
```

#### Obtenir un lien
```bash
# Afficher les détails
al link get #github
al link get documentation

# Copier l'URL dans le presse-papier
al link get #api --cp
allink get rest --cp
```

#### Éditer un lien
```bash
# Changer l'URL
al link edit #github -u https://github.com/newuser/repo

# Ajouter des keywords
al link edit #docs --add_keyword tutorial|tuto
al link edit #docs -ak new_keyword

# Réinitialiser les keywords
al link edit #api --reset_keyword rest|api
al link edit #api -rk

# Supprimer tous les keywords
al link edit #api -rk ""
```

#### Supprimer un lien
```bash
al link remove #github
allink remove api
```

## Options communes

- `-t, --target <project>` : Spécifier un projet cible différent du projet courant
- `-c, --chiffre` : Chiffrer une note (commande note uniquement)
- `-b, --body <content>` : Spécifier le contenu directement sans ouvrir l'éditeur
- `--cp` : Copier dans le presse-papier au lieu d'afficher

## Structure des données

### Répertoire global : `~/.al_global/`
- `projects` : Mapping des projets avec leurs raccourcis et chemins
- `config` : Configuration globale (longueur de prévisualisation, etc.)

### Répertoire local : `<projet>/.al_local/`
- `notes/` : Notes du projet (JSON)
- `links/` : Liens du projet (JSON)

## Notes chiffrées

Les notes chiffrées utilisent AES-GCM avec un mot de passe. Le système inclut une chaîne de vérification pour détecter un mot de passe incorrect lors du déchiffrement.

## Suggestions de noms similaires

Quand un nom de note ou de lien n'est pas trouvé, la CLI suggère des noms similaires basés sur la distance de Levenshtein.

## Développement

### Compiler
```bash
make build
```

### Nettoyer
```bash
make clean
```

### Tests
```bash
make test
```

## License

MIT
