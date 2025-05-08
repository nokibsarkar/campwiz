SFTP_HOST=login.toolforge.org
SFTP_USER=nokibsarkar
SFTP_PRIVATE_KEY=~/.ssh/id_rsa
SFTP_PORT=22
SFTP_DIR=/data/project/backup-bot/backups/CampWiz-NXT
DB_HOST=localhost
DB_NAME=campwiz
DB_USER=campwiz
DB_PASS=$(cat .env | grep DB_PASS | cut -d '=' -f2)
backup_file_prefix=campwiz-backup-$( date +%Y-%m-%d_%H-%M-%S )
backup_file_raw=${backup_file_prefix}.sql
backup_file_gz=${backup_file_prefix}.sql.tar.gz
echo "Creating backup of the database..."
mysqldump --host=$DB_HOST --user=$DB_USER --password=$DB_PASS $DB_NAME > $backup_file_raw
echo "Compressing the backup file..."
tar -cvzf $backup_file_gz $backup_file_raw
echo "Removing the raw backup file..."
rm $backup_file_raw
echo "Backup created: $backup_file_gz"
echo "Backup file size: $(du -h $backup_file_gz | cut -f1)"
echo "Backup file location: $(pwd)/$backup_file_gz"
echo "Transfer the backup file to a secure location."
sftp -i $SFTP_PRIVATE_KEY -P $SFTP_PORT $SFTP_USER@$SFTP_HOST <<EOF
cd $SFTP_DIR
put $backup_file_gz
bye
EOF
echo "Backup file transferred to SFTP server: $SFTP_HOST"