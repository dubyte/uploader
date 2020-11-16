# uploader 

It is a simple web app to upload file with an html form and a web server. The idea is to be used with a reverse_proxy that add the https, like caddy or nginx.

## Deployment for ubuntu arm64

```bash
    # build 
    GOARCH=arm64 go build .

    # copy to your server
    scp uploader user@host:~
        
    ssh user@host
    cd
    sudo chmod 755 uploader
    sudo chown root:root uploader
    sudo mv uploader /usr/local/bin

    # create a systemd service 
    sudo touch /etc/systemd/system/uploader.service
    sudo chmod 664 /etc/systemd/system/uploader.service

    # copy all until EOT
    sudo tee -a /etc/systemd/system/uploader.service > /dev/null <<EOT
    [Unit]
    Description=A web app to upload files to the server

    [Service]
    PIDFile=/run/uploader.pid
    ExecStart=/usr/local/bin/uploader 
    User=
    Group=

    [Install]
    WantedBy=multi-user.target
    EOT

    # test service run
    sudo systemctl start uploader
    sudo systemctl status uploader

    # if everything looks ok
    sudo systemctl enable uploader

```

## HTTPS
Remember if you are not using https your files could be intercepted. I am using caddy

```bash
    # you have to have a domain name so caddy can create a self signed certificate for https
    
    # for internal lan autohttps use:
    # echo 'tls internal' | sudo tee -a /etc/caddy/Caddyfile > /dev/null
    
    # append to your Caddyfile
    echo 'reverse_proxy /upload localhost:8080' | sudo tee -a /etc/caddy/Caddyfile > /dev/null
    
    sudo systemctl restart caddy
```
## Notes
