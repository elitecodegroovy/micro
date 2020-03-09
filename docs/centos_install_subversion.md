## Step 1: Update the system
```
sudo yum update
sudo shutdown -r now
```

After the reboot, log into the system again using the same sudo user.

## Step 2: Install Apache
Install Apache using YUM:

```
sudo yum install httpd
```

Remove Apache's default welcome page:

```
sudo sed -i 's/^/#&/g' /etc/httpd/conf.d/welcome.conf

```

Prevent Apache from displaying files in the "/var/www/html" directory:

```
sudo sed -i "s/Options Indexes FollowSymLinks/Options FollowSymLinks/" /etc/httpd/conf/httpd.conf
```

## Step 3: Install SVN and the mod_dav_svn module
To make SVN work with Apache, you need to install an Apache module "moddavsvn" in addition to SVN:

```
sudo yum install subversion mod_dav_svn

```

## Step 4: Configure SVN

1) Modify the SVN configuration file

```
sudo vi /etc/httpd/conf.modules.d/10-subversion.conf
```

The file should read:

```
LoadModule dav_svn_module     modules/mod_dav_svn.so
LoadModule authz_svn_module   modules/mod_authz_svn.so
LoadModule dontdothat_module  modules/mod_dontdothat.so
```

Append the following segment:
```
<Location /svn>
DAV svn
SVNParentPath /svn
AuthName "SVN Repos"
AuthType Basic
AuthUserFile /etc/svn/svn-auth
AuthzSVNAccessFile /svn/authz
Require valid-user
</Location>
```

Save and quit:
```
:wq!
```

Note: In this configuration, we specified an HTTP access authentication file "/etc/svn/svn-auth" and a user permission control file "/svn/authz". Both of them will be created later.

2) Create an SVN repo
```
sudo mkdir /svn
cd /svn
sudo svnadmin create repo1
sudo chown -R apache:apache repo1

```

3) Setup SVN user accounts
Use the following commands to create an HTTP access authentication file "/svn/svn-auth" and an SVN user account "user001":
```
sudo mkdir /etc/svn
sudo htpasswd -cm /etc/svn/svn-auth user001
sudo chown root:apache /etc/svn/svn-auth
sudo chmod 640 /etc/svn/svn-auth
```

If you want to create more SVN user accounts, refer to the following commands:

```
sudo htpasswd -m /etc/svn/svn-auth user002
sudo htpasswd -m /etc/svn/svn-auth user003

```
Warning: Do not use the "-c" flag from now on, or you will rebuild the authentication file and erase all the user accounts you setup earlier.

4) Setup permissions for users
```
sudo cp /svn/repo1/conf/authz /svn/authz
sudo vi /svn/authz

```

Assume that:

User "user001" is the administrator.
User "user002" is a qualified user who owns read and write permissions to the SVN repo "repo1".
User "user003" is a trainee who can only read the contents of SVN repo "repo1".
Then you can modify the settings within as below:
```
[groups]
admin=user001
repo1_user=user002
repo1_trainee=user003

[/]
@admin=rw

[repo1:/]
@repo1_user=rw
@repo1_trainee=r
Save and quit:

:wq!

```
## Step 5: Start Apache and modify firewall rules

Start Apache:
```
sudo systemctl start httpd.service
sudo systemctl enable httpd.service
```

Open the HTTP service port:

```
sudo firewall-cmd --zone=public --permanent --add-service=http
sudo firewall-cmd --reload

```

Finally, use the following path to access the SVN repo "repo1" on your server from an SVN client:

```
http://<your-server-ip>/svn/repo1/
```


That concludes our tutorial. Thank you for reading.