pipeline {
    agent any

    parameters {
        booleanParam(
            name: 'PUSH_IMAGE',
            defaultValue: false,
            description: 'Push image to configured registry'
        )
    }

    environment {
        APP_NAME = 'shift-devops-app'
        LOCAL_IMAGE = 'shift-devops'
        REGISTRY = 'ghcr.io'
        REGISTRY_IMAGE = 'ghcr.io/hrthvln/shift-devops-test'
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

        stage('Build Image') {
            steps {
                script {
                    env.VERSION = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()

                    env.IMAGE_TAG = "${LOCAL_IMAGE}:${VERSION}"

                    sh '''
                        docker build \
                          --build-arg VERSION="$VERSION" \
                          -t "$IMAGE_TAG" \
                          .
                    '''
                }
            }
        }

        stage('Push Image') {
            when {
                expression {
                    return params.PUSH_IMAGE
                }
            }
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: 'registry-credentials',
                        usernameVariable: 'REGISTRY_USER',
                        passwordVariable: 'REGISTRY_PASS'
                    )
                ]) {
                    sh '''
                        echo "$REGISTRY_PASS" | docker login "$REGISTRY" \
                          --username "$REGISTRY_USER" \
                          --password-stdin

                        docker tag "$IMAGE_TAG" "${REGISTRY_IMAGE}:${VERSION}"
                        docker push "${REGISTRY_IMAGE}:${VERSION}"

                        docker logout "$REGISTRY"
                    '''
                }
            }
        }

        stage('Prepare Hotfix Binary') {
            steps {
                sh '''
                    rm -f hotfix-app

                    docker create \
                      --name shift-devops-extract-"$BUILD_NUMBER" \
                      "$IMAGE_TAG"

                    docker cp \
                      shift-devops-extract-"$BUILD_NUMBER":/app \
                      hotfix-app

                    docker rm \
                      shift-devops-extract-"$BUILD_NUMBER"
                '''
            }
        }

        stage('Deploy') {
            steps {
                sh '''
                    set -e

                    CONTAINER="$APP_NAME"
                    BACKUP="backup-app"

                    if ! docker inspect "$CONTAINER" >/dev/null 2>&1; then
                        echo "ERROR: target container $CONTAINER does not exist."
                        exit 1
                    fi

                    docker cp "$CONTAINER:/app" "$BACKUP"

                    docker cp hotfix-app "$CONTAINER:/app.new"

                    docker exec "$CONTAINER" \
                      sh -c "chmod +x /app.new && mv /app.new /app"

                    if ! docker restart "$CONTAINER"; then
                        echo "Restart failed. Rolling back."
                        docker cp "$BACKUP" "$CONTAINER:/app.rollback"
                        docker exec "$CONTAINER" \
                          sh -c "chmod +x /app.rollback && mv /app.rollback /app"
                        docker restart "$CONTAINER" || true
                        exit 1
                    fi

                    success=false

                    for i in 1 2 3 4 5 6 7 8 9 10; do
                        response="$(curl -s http://host.docker.internal:8080/ || true)"

                        if echo "$response" | grep -q "version=$VERSION"; then
                            success=true
                            echo "Deployment verification passed:"
                            echo "$response"
                            break
                        fi

                        sleep 1
                    done

                    if [ "$success" != "true" ]; then
                        echo "Health/version check failed. Rolling back."

                        docker cp "$BACKUP" "$CONTAINER:/app.rollback"
                        docker exec "$CONTAINER" \
                          sh -c "chmod +x /app.rollback && mv /app.rollback /app"

                        docker restart "$CONTAINER" || true

                        exit 1
                    fi
                '''
            }
        }
    }

    post {
        always {
            sh '''
                rm -f hotfix-app backup-app || true
            '''
        }
    }
}