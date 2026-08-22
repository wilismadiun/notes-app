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
        stage('test') {
            steps {
                sh './jenkins/scripts/test.sh'
            }
        }
    }
}