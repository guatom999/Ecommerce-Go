pipeline {
    agent any

    environment {
        IMAGE_NAME = "ecommerce-go"
        IMAGE_TAG  = "latest"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Test') {
            steps {
                sh 'go test ./...'
            }
        }

        // stage('Build') {
        //     steps {
        //         sh 'docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .'
        //     }
        // }
    }

    post {
        success {
            echo '✅ Pipeline succeeded!'
        }
        failure {
            echo '❌ Pipeline failed!'
        }
    }
}
