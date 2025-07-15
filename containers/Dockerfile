    FROM ubuntu:22.04

    ENV container docker

    # Install only essential packages with no extras
    RUN apt-get update && \
        apt-get install -y --no-install-recommends \
            systemd \
            systemd-sysv \
            dbus && \
        apt-get clean 
    # Install OpenSSH Server
    RUN apt-get update && \
        apt-get install -y openssh-server && \
        mkdir -p /var/run/sshd

    RUN apt-get update && apt-get install -y openssh-client

    # Generate SSH host keys (if not already present)
    RUN ssh-keygen -A

    # Optional: set root password (insecure, change for production!)
    # RUN echo 'root:root' | chpasswd

    # Allow root login (optional & insecure for production)
    RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

    

    #  insert the public key into the authorized_keys file
    COPY ./containers/id_rsa.pub /root/.ssh/authorized_keys

    # Expose SSH port
    EXPOSE 22

    VOLUME [ "/sys/fs/cgroup" ]
    STOPSIGNAL SIGRTMIN+3

    # Enable sshd service to start on boot
    RUN systemctl enable ssh

    CMD ["/sbin/init"]
    
    # Start the SSH service when the container runs
    # CMD ["/usr/sbin/sshd", "-D"]



 


   