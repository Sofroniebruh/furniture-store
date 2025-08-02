package org.example;

public class EmailPOJO
{
    private String email;
    private String body;
    private String subject;

    public EmailPOJO()
    {
    }

    public EmailPOJO(String email, String body, String subject)
    {
        this.email = email;
        this.body = body;
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

    public String getBody()
    {
        return body;
    }

    public void setBody(String body)
    {
        this.body = body;
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
