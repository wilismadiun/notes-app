pipeline {
    agent {
        docker {
            image 'golang:1.26'
            args '-p 3000:3000'
        }
    }
    environment {
        // Mencegah error permission denied saat Go membuat cache build
        GOCACHE = '/tmp/go-cache'
        GOPATH = '/tmp/go-path'
    }
    stages {
        stage('Restore Depedencies') {
            steps {
                sh 'go mod download'
            }
        }
        stage('setup environment') {
            steps {
                sh '''
                    echo "DB_HOST=host.docker.internal" > .test.env
                    echo "DB_PORT=5432" >> .test.env
                    echo "DB_USER=developers" >> .test.env
                    echo "DB_PASSWORD=supersecretpassword" >> .test.env
                    echo "DB_NAME=notesapp_test" >> .test.env
                    echo "DB_SSLMODE=disable" >> .test.env
                    echo "JWT_SECRET=0ff69fbff76d8dd13513a51c2e6db40584b7aa27c8cca2eaf6427ede666ab6b8" >> .test.env
                '''
            }
        }
        stage('test') {
            steps {
                sh 'go test -v ./...'
            }
        }
    }
}