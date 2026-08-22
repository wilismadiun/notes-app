pipeline {
    agent {
        docker {
            image 'golang:1.26'
            args '-p 3000:3000'
        }
    }
    stages {
        stage('Restore Depedencies') {
            steps {
                sh 'go mod download'
            }
        }
    }
}