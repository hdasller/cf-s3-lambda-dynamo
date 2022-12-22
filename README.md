# cf-s3-lambda-dynamodb
Projeto sobre o aprendizado adquirido com as ferramentas Cloud Formation, S3, lambda e dynamodb

## Rodando o projeto
Para realizar o deploy na aws utilizando o cloudformation, deve executar o seguinte comando:

`./infrastructure/apply  <nome-stack> <caminho-yaml-cloud-formation> <caminho-parametros-json>`

## Buildando o projeto

`docker build -t <nome-imagem> . --no-cache -f infrastructure/Dockerfile`

`docker tag <nome-imagem> <endereço-repo-imagem>`

`docker push  <endereço-repo-imagem>`

## Melhorias futuras

- Adicionar resource do api gateway no cloudformation 
- Adicionar resource do codepipeline no cloudformation
- Adicionar secrets


## Video do funcionamento
https://www.loom.com/share/a0c3304617b64936b25c351a766d7723



