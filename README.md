# transactions-processor

El servicio realiza el procesamiento de un archivo alojado en un bucket de S3 de tipo CSV.

Luego almacena los datos en un registro de dynamoDb el cual dispara un envio de reporte via email.


Stack/Obs:
- Golang
- AWS
- Proteccion de branches main y develop
- Uso de conventional commits

WIP