# Fanouter
<<<<<<< HEAD
=======
[![Build Status](https://travis-ci.com/shipa988/fanouter.svg?branch=master)](https://travis-ci.com/shipa988/fanouter)
[![Go Report Card](https://goreportcard.com/badge/github.com/shipa988/fanouter)](https://goreportcard.com/report/github.com/shipa988/fanouter)
>>>>>>> 0ed0c24dc0cdaa54b3acc86802ceb6f47925181e

service for qps limiting and fanouting

test job for Gadsbee

-Необходимо сделать веб-сервер, который будет принимать входящие запросы и обрабатывать данные исходя из URL параметров. Формат и данные параметров формируются произвольно, не принципиально.

-Необходимо сформировать любые 10 внешних URL, на которые в один момент времени после приема входящего запроса будут отправляться запросы

-Необходимо организовать возможность лимитировать количество исходящих запросов по любому из внешних URL

-структуры можно хранить как внутри программы, так и загружать из внешнего JSON файла.

-На выходе должны получить возможность менять значение исходящих QPS для каждого внешнего URL.

-покрыть весь код тестами, чтобы убедиться, что лимит по QPS работает
