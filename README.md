# go-luminaria
Projeto que controla um relay que controla uma luminaria usando Gobot (Go)

## Compilando
```
env GOOS=linux GOARCH=arm GOARM=6 go build -o controleluminaria-arm-v6
env GOOS=linux GOARCH=arm GOARM=7 go build -o controleluminaria-arm-v7
```

## Copiando ao Raspberry Pi
```
scp controleluminaria-arm-v6 pi@192.168.0.112:/home/pi/controleluminaria-arm-v6
```

## Iniciando a aplicação no Raspberry Pi
```
sudo nohup ./controleluminaria-arm-v6 > controle.log 2>&1 &
``` 

