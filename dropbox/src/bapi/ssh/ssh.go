package ssh 

import (
    "fmt"
    "golang.org/x/crypto/ssh"
    "log"
    "time"
)

func connect(user, password, host string, port int) (*ssh.Session, error) {
    var (
        auth         []ssh.AuthMethod
        addr         string
        clientConfig *ssh.ClientConfig
        client       *ssh.Client
        session      *ssh.Session
        err          error
    )
    auth = make([]ssh.AuthMethod, 0)
    auth = append(auth, ssh.Password(password))
    clientConfig = &ssh.ClientConfig{
        User:    user,
        Auth:    auth,
        Timeout: 30 * time.Second,
    }
    addr = fmt.Sprintf("%s:%d", host, port)
    if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
        return nil, err
    }
    if session, err = client.NewSession(); err != nil {
        return nil, err
    }
    return session, nil
}
func Ssh(ip, cmd string) (string, error) {
    session, err := connect("xxx", "xxx", ip, 22)
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()
    //session.Stdout = os.Stdout
    //session.Stderr = os.Stderr
    //session.Run(cmd)
    b, err := session.Output(cmd)
    /*
        if err != nil {
            log.Fatalf("failed to execute: %s", err)
        }*/
    fmt.Println(string(b))
    return string(b), err
}

