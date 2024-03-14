# transactions-processor

El servicio realiza el procesamiento de un archivo alojado en un bucket de S3 de tipo CSV.

Luego almacena los datos en un registro de dynamoDb el cual dispara un envio de reporte via email.

El flujo de funcionamiento es el siguiente>

1) Se realiza un GET a la siguiente URL ``https://sov0g958ra.execute-api.us-east-1.amazonaws.com/process``
2) Este endpoint tiene integrada una lambda (transactions-processing) que:
   - realiza la obtencion del archivo CSV del bucket de S3  
   - Procesa las transacciones halladas en el file 
   - Persiste la informacion de la cuenta en un registro en una tabla de dynamoDb
3) Existe otra lambda (account-balance-email-sender) integrada para escuchar actualizaciones o inserts de registros de dynamoDB y hacer el envio del email utilizando AWS SES


Stack/Obs:
- Golang
- AWS
- Proteccion de branches main y develop
- Uso de conventional commits
- Inclusion basica de tests unitarios aplicando TDT solo un handler solamente

