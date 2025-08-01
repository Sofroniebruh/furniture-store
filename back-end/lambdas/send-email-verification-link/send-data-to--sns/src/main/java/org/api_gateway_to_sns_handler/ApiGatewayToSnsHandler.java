package org.api_gateway_to_sns_handler;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyRequestEvent;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyResponseEvent;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import software.amazon.awssdk.services.sns.SnsClient;
import software.amazon.awssdk.services.sns.model.PublishRequest;

import java.util.Map;

public class ApiGatewayToSnsHandler implements RequestHandler<APIGatewayProxyRequestEvent, APIGatewayProxyResponseEvent>
{
    private static final SnsClient snsClient = SnsClient.create();
    private static final ObjectMapper objectMapper = new ObjectMapper();

    private static final String TOPIC_ARN = System.getenv("SNS_TOPIC_ARN");

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
            Map<String, Object> json = objectMapper.readValue(request.getBody(), Map.class);
            String message = "Email: " + json.get("email") + "\n" +
                    "Message: " + json.get("message") + "\n";
            PublishRequest publishRequest = PublishRequest.builder()
                    .topicArn(TOPIC_ARN)
                    .message(message)
                    .build();
            snsClient.publish(publishRequest);

            String value = convertValueToJson("message: ", "Email was send successfully");

            return new APIGatewayProxyResponseEvent()
                    .withStatusCode(200)
                    .withBody(value);
        }
        catch (Exception e)
        {
            try
            {
                String value = convertValueToJson("error: ", e.getMessage());

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
