package org.example;

public class EmailPOJO
{
    private String email;
    private String messageBody;
    private String subject;

    public EmailPOJO()
    {
    }

    public EmailPOJO(String email, String messageBody, String subject)
    {
        this.email = email;
        this.messageBody = messageBody;
        this.subject = subject;
    }

    public String getEmail()
    {
        return email;
    }

    public void setEmail(String email)
    {
        this.email = email;
    }

    public String getMessageBody()
    {
        return messageBody;
    }

    public void setMessageBody(String messageBody)
    {
        this.messageBody = messageBody;
    }

    public String getSubject()
    {
        return subject;
    }

    public void setSubject(String subject)
    {
        this.subject = subject;
    }
}
