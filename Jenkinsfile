pipeline {
    agent {
        dokcer {
            image 'golang:1.26'
            args '-p 3000:3000'
        }
    }
    stages {
        stage('Restore Depedencies') {
            sh 'go mod download'
        }
    }
}