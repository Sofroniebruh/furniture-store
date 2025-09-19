package org.example;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyRequestEvent;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyResponseEvent;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import software.amazon.awssdk.services.ses.SesClient;
import software.amazon.awssdk.services.ses.model.*;
import java.util.Map;

public class ApiGatewayToSesHandler implements RequestHandler<APIGatewayProxyRequestEvent, APIGatewayProxyResponseEvent>
{
    private static final ObjectMapper objectMapper = new ObjectMapper();
    private static final SesClient sesClient = SesClient.create();
    private static final String verifiedEmail = System.getenv("VERIFIED_EMAIL");

    public boolean sendEmail(String receiverEmail, String subject, String emailBody)
    {
        System.out.println("I am here");
        System.out.println("Email body: " + emailBody);
        
        if (subject == null) subject = "";
        if (emailBody == null) emailBody = "";
        
        try
        {
            Destination destination = Destination.builder()
                    .toAddresses(receiverEmail)
                    .build();
            Content subjectContent = Content.builder()
                    .data(subject)
                    .build();
            Content bodyContent = Content.builder()
                    .data(emailBody)
                    .build();

            Body body = Body.builder()
                    .text(bodyContent)
                    .build();
            Message message = Message.builder()
                    .subject(subjectContent)
                    .body(body)
                    .build();
            SendEmailRequest emailRequest = SendEmailRequest.builder()
                    .destination(destination)
                    .message(message)
                    .source(verifiedEmail)
                    .configurationSetName(System.getenv("SES_CONFIGURATION_SET"))
                    .build();

            sesClient.sendEmail(emailRequest);

            return true;
        }
        catch (Exception e)
        {
            System.out.println(e);
            return false;
        }
    }

    public String convertValueToJson(String value1, String value2) throws JsonProcessingException
    {
        Map<String, Object> value = Map.of(
                value1, value2
        );

        return objectMapper.writeValueAsString(value);
    }

    @Override
    public APIGatewayProxyResponseEvent handleRequest(APIGatewayProxyRequestEvent request, Context context)
    {
        try
        {
            System.out.println("Raw request body: " + request.getBody());
            System.out.println("Request headers: " + request.getHeaders());
            
            String requestBody = request.getBody();
            // Handle case where body might be double-encoded as a JSON string
            if (requestBody.startsWith("\"") && requestBody.endsWith("\"")) {
                requestBody = objectMapper.readValue(requestBody, String.class);
            }
            
            EmailPOJO emailPOJO = objectMapper.readValue(requestBody, EmailPOJO.class);
            String receiverEmail = emailPOJO.getEmail();
            String body = emailPOJO.getMessageBody();
            String emailBody;
            String subject = emailPOJO.getSubject();

            System.out.println("Your email POJO: " + emailPOJO);
            System.out.println("Your email body: " + body);
            
            if (receiverEmail == null || receiverEmail.trim().isEmpty()) {
                throw new CustomException("Receiver email is required");
            }
            if (subject == null || subject.trim().isEmpty()) {
                throw new CustomException("Subject is required");
            }

            System.out.println("Subject for comparison: '" + subject + "'");
            
            if (subject.toLowerCase().contains("verify"))
            {
                if (body == null) body = "";
                emailBody = "Your verification code: " + body + ". The code is valid for 2 minutes.";
            }
            else if (subject.toLowerCase().contains("reset") || subject.toLowerCase().contains("password"))
            {
                if (body == null) body = "";
                emailBody = "Please go to the following url to reset your password: " + body;
            }
            else
            {
                emailBody = body != null ? body : "";
            }

            if (!sendEmail(receiverEmail, subject, emailBody)) throw new CustomException("Something went wrong");

            String value = convertValueToJson("message", "Email was sent successfully");

            return new APIGatewayProxyResponseEvent()
                    .withStatusCode(200)
                    .withBody(value);
        }
        catch (Exception e)
        {
            try
            {
                String value = convertValueToJson("error", e.getMessage());

                return new APIGatewayProxyResponseEvent()
                        .withStatusCode(500)
                        .withBody(value);
            }
            catch (JsonProcessingException jpe)
            {
                return new APIGatewayProxyResponseEvent()
                        .withStatusCode(500)
                        .withBody(e.getMessage());
            }
        }
    }
}
