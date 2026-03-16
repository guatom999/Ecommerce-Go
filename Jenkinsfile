pipeline {
    agent any

    tools {
        go 'go-1.23'
    }

    environment {
        IMAGE_NAME        = "ecommerce-go"
        IMAGE_TAG         = "latest"
        COVERAGE_THRESHOLD = 70
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Test Two') {
            steps {
                script{
                    echo "test"
                }
            }
        }

        stage('Test') {
            steps {
                script {
                    // รัน test พร้อมเก็บ coverage และเก็บ output
                    def output = sh(
                        script: 'go test -coverprofile=coverage.out ./... 2>&1',
                        returnStdout: true
                    ).trim()

                    echo output

                    // หาว่ามี package ไหนที่ coverage >= COVERAGE_THRESHOLD ไหม
                    def threshold = env.COVERAGE_THRESHOLD.toDouble()
                    def passed = output.readLines().any { line ->
                        def matcher = line =~ /coverage: ([\d.]+)% of statements/
                        if (matcher) {
                            return matcher[0][1].toDouble() >= threshold
                        }
                        return false
                    }

                    if (!passed) {
                        error "❌ ไม่มี layer ไหนที่มี coverage >= ${threshold}% — pipeline ล้มเหลว!"
                    }

                    echo "✅ มีอย่างน้อย 1 layer ที่ coverage >= ${threshold}%"
                }
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
