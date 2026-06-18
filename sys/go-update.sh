#!/usr/bin/bash

function GO_UPD() {
  # Récupération de la version en ligne
  # L'option -k a été retirée de curl pour forcer la validation du certificat TLS (sécurité)
  local GO_VER_RAW=$(curl -s "https://go.dev/VERSION?m=text")
  local GO_VER=$(echo "$GO_VER_RAW" | awk 'NR==1')

  if [[ -z "$GO_VER" ]]; then
    echo "Erreur : Impossible de récupérer la dernière version de Go en ligne."
    return 1
  fi

  # Récupération de la version locale
  local LOCAL_GO_VER="none"
  if command -v go &> /dev/null; then
    # Extrait par exemple "go1.22.1" de la commande "go version"
    LOCAL_GO_VER=$(go version | awk '{print $3}')
  fi

  # Comparaison des versions
  if [[ "$LOCAL_GO_VER" == "$GO_VER" ]]; then
    echo " + Go est déjà à jour (Version actuelle : ${LOCAL_GO_VER})"
    return 0
  fi

  echo " + Nouvelle version détectée : ${LOCAL_GO_VER} -> ${GO_VER}"
  echo " + Téléchargement de ${GO_VER} linux amd64..."
  
  local TAR_FILE="${GO_VER}.linux-amd64.tar.gz"
  wget -q "https://golang.org/dl/${TAR_FILE}"

  echo " + Installation en cours..."
  # Pratique sécurisée officielle : supprimer l'ancienne installation avant d'extraire la nouvelle
  # Cela évite que des fichiers obsolètes ne corrompent le système
  sudo rm -rf /usr/local/go
  sudo tar xf "${TAR_FILE}" -C /usr/local

  # Nettoyage du fichier téléchargé
  rm "${TAR_FILE}"
  go clean -cache

  echo " + Mise à jour terminée."
  
  # Message utilisateur.
  echo " + N'oubliez pas d'exécuter 'source ~/.bashrc' si ce n'est pas déjà fait."
}

GO_UPD

