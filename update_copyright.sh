#!/bin/sh
#
# script de remplacement des Copyright de headers
#
# Copyright jean-francois PHILIPPE  2011-2018
#

#
# Met a jour un type de fichier
# Parametres
#  $1 Extension de fichier
#  $2 la nouvelle annee
#
update_files_type() {
echo "Fichier de type $1  annee $2"
echo "    Copyright sur 2 ans"
    # MAJ des copyright sur 2 annees
    find . -type f -iname "$1" -exec sed -i "s,\(Copyright .* PHILIPPE *[0-9]\{4\}\)\(-[0-9]\{4\}\),\1-$2,g" {} \;
    # MAJ des copyright sur 1 annee
echo "    Copyright sur 1 ans 1/2"
    find . -type f -iname "$1" -exec sed -i "s,\(Copyright .* PHILIPPE *[0-9]\{4\}\)\([^-]\),\1-$2\2,g" {} \;
    # MAJ des copyright sur 1 annee et l annee en fin de ligne
echo "    Copyright sur 1 ans 2/2"
    find . -type f -iname "$1" -exec sed -i "s,\(Copyright .* PHILIPPE *[0-9]\{4\}\)$,\1-$2,g" {} \;
}

YEAR=$(date +"%Y")
update_files_type '*.go'  "$YEAR"
update_files_type '*.sh' "$YEAR"
update_files_type '*.sql' "$YEAR"
update_files_type 'LICENSE' "$YEAR"
