#!/bin/bash

# Токен для авторизации (замени на свой после выполнения /login)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InBsYXllcjEiLCJleHAiOjE3NDA4MTQzMTV9.2OIWeGPMK07ppPFypAFizgXRDTAySrXYxgqm4Yt-wSQ"  # Вставь свой токен

# Количество запросов
REQUESTS=100

echo "Запускаем тест: $REQUESTS запросов на /submit-score"

for ((i=1; i<=REQUESTS; i++))
do
    SCORE=$((RANDOM % 1000 + 1))  # Случайный счёт от 1 до 1000
    curl -X POST -H "Content-Type: application/json" \
         -H "Authorization: Bearer $TOKEN" \
         -d "{\"game\":\"game1\",\"score\":$SCORE}" \
         http://localhost:8080/submit-score > /dev/null 2>&1 &
    echo "Отправлен запрос $i с очками $SCORE"
done

echo "Тест завершён. Проверьте логи сервера."