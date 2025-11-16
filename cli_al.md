Ceci est la description d'une CLI nommé 'al' et programmé en Golang.
Son but est de me faciliter la gestion des projets de mes clients en ligne de commande.
Pouvoir stocker des notes liées au projet, des raccourcis de commande, des liens https, etc... 
Elle sera compatible sur Ubuntu pour commencer.

De la même manière que GIT, quand j'intiialise un projet Al, il crée un dossier .al_local dans le répertoire courante mais également, il stocke de manière centralisée 
dans /home/$user/.al_global un dictionnaire de correspondance "nom du dossier" -> Path complet. De tel manière que quand je fais un 'al go $nom' il change le répertoire courant du shell actuel par le répertoire du projet en question.

Petite précision importante, chaque action de la cli aura sa propre CLI au sens ou 
je peux la déclencher à la fois en faisant "al action" mais également "alaction" 
al go marche aussi avec algo
al note marche aussi avec alnote
Etc.... SAUF al install et al update



#1 - 'al install'

Quand on fait appelle à la commande install, probablement depuis le répertoire de build du projet,
la cli doit, 1, créer un répertoire '.al_global' dans /home/$user/ si il n'existe pas ET deux, 
il doit prendre ses propres binaires et les copier dans le repertoire bin
qui permettra d'accéder à toutes les commandes 'al' dans le shell, n'importe où.

Dans le repertoire .al_global, on va retrouver un fichier 'config' et un fichier 'projects' 
Le fichier config stockera des paramètres de configuration et le fichier project, une correspondance "[shortcuts]" -> "path to projects"

#1.5 'al update'
Comme al install, lancée depuis le répertoire de build du projet, l'update ne va juste replacer les binaires déjà présent dans /bin et ajouter les manquants/nouveaux.

#2 - al init 'shortcut|...' ou alinit 

La commande init décide que le répertoire courant sera un projet. Il crée donc un dossier .al dans le répertoire courant et crée une nouvelle entrée dans le fichier projects du dossier global.
quand on initialise on peut spécifier des shortcuts ou keywords associés au projet, qui permettront d'y faire référence.
Par défaut, il prend systèmatiquement le nom du dossier courant comme name/shortcuts, et ensuite, tous les noms qui sont séparés par un pipe |
Par exemple, je suis dans le repertoirre ./foo et je fais un al init bar|mad|tes l'ensemble des raccourcis à sauvegarder son foo, bar, mad, tes

#3 - al go 'name' ou algo
Une commande qui permet directement d'aller dans le répertoire du projet dont le nom ou les shortcuts correspond à 'name'
par exemple, j'ai init le répertoire /home/alex/testfoo avec les shortcuts test et foo alors je peux faire n'importe où dans le shell :
al go testfoo
al go test
al go foo
et je me retrouverais dans /home/alex/testfoo.

#4 - al note [action] #name_note [arguments] ou alnote  
La commande 'note' permet de gérer les notes pour un projet.
Les notes sont juste des bouts de textes que je souhaite conserver, modifier et facilement accéder.

Certaines notes seront sensibles et donc chiffrée. (c'est à dire qu'elles ne seront pas stockées en clair et qu'une clé de déchiffrement (simple password) sera demandée à l'ajout et à l'édit). Quand elles sont chiffrées, le programme rajoute une chaine de caractère dedans pour savoir si quand elle déchiffre avec le password donné par l'user, le password est le bon.

les sous-actions sont get|add|edit|remove|list
Par défaut, les notes ciblées concernes le projet courant (celui du répertoire courant) mais si il y à l'argument -target (ou -t) avec un nom ou un shortcut spécifié juste après, c'est dans le projet cible que ça s'execute.
Les sous-actions  :
### Lister :
    'al note list' => Liste les notes concernant le projet courant
    'al note list -t/--target foo' => Liste les notes concernant le projet à destination de 'foo'
    La liste c'est un tableau :
    Nom -> Date -> Apercu du contenu 
    avec triée par Nom, et l'apercu du contenu étant les 60 premières caractères (sachant que 60 peut changer)
    Pour les notes chiffrées dans l'apercu tu n'afficheras juste "**chiffrée**"
### Ajouter :
   'al note add #name' => Ajoute une note ayant le nom #name pour le projet courant et ouvre un vim sur la note.
   'al note add #name -t/--target foo' => Ajoute une note ayant le nom #name pour le projet foo et ouvre le vim de la note.
   'al note add #name -c/--chiffre => Ajoute une note chiffrée ayant le nom #name, il demande d'abord le password en ligne de commande puis ouvre un vim sur la note.
   'al note add #name -b/--body' => ajoute une note ayant le nom #name et tout ce qui vient après l'argument body est automatiquement ajoutée à la note sans ouverture de vim. c'est à l'utilisateur de s'assurer de le mettre en dernier sinon les autres paramètres ne sont pas pris en compte.

   évidemmment, tous les arguments sont cumulables.
   on peut donc faire 'al add #name -t foo -c -b "blablablabla"

### Obtenir : 
   'al note get #name' => affiche dans le shell le contenu de la note #name, demande le password avant si elle est chiffrée.
   'al note get #name -t/--target foo' précise un projet distant.
   'al note get #name -cp => N'affiche pas mais copie dans le presse papier le contenu de la note. Demande le password avant si elle est chiffrée.

### Edit : 
  'al note edit #name' => Ouvre un vim sur la note à fin de la modifier. Demander le password avant si elle est chiffrée.
  'al note edit #name -t/--target foo' précise un projet distant
  'al note edit #name -b/--body LOREM IPSUM ... => écrase l'ensemble du corps de la note par le body, demande le password avant si chiffrée.

### Supprimer :
  'al note remove #name' => Supprime la note après avoir demandé confirmation.
  'al note remove #name -t/--target foo => Pareil mais avec projet distant

Précisions pour l'ensemble des sous-actions si le #name correspond pas, le programme renvoi les noms de notes qui pourraient correspondre basées sur un calcul de distance.

#5 -  al link ou allink 
Permet de gérer des liens https. Les liens auront des noms ET des keywords/shortcuts pour pouvoir s'y référer de manière différente.
Comme pour les autres commandes, à la base c'est le projet du répertoire courant qui est visée et sinon, -t/--target pour désigner un autre projet.
Les sous-actions sont list/add/get/edit 

### Lister 
    'al link list'
    donne la liste au format Nom -> Link -> Keywords

### Ajouter 
    'al link add #name -u/--url url -k/--keywords foo|bar|maz' 

### Obtenir 
    'al link get #name/#keyword' affiche le link ET les keywords
    'al link get #name/#keyword -cp' le paramètre cp copie colle direction le link dans le presse papier

### edit 
    'al link edit #name/#keyword -u/--url url pour changer l'url
    'al link edit #name/#keyword --add_keyword/-ak foo|bar pour ajouter des keywords.
    'al link edit #name/#keyword --reset_keyword/-rk foo|bar pour écraser les keywords existant. (si rien après -rk alors écraser par vide)
### remove
    'al link remove #name/#keyword' supprime après demande de confirmation.

Précisions pour l'ensemble des sous-actions si le #name correspond pas, le programme renvoi les noms de notes qui pourraient correspondre basées sur un calcul de distance.




