using FluentValidation;
using Microsoft.AspNetCore.Diagnostics;
using Microsoft.AspNetCore.Mvc;
using AccountService.API.Exceptions;

namespace AccountService.API.Handlers;

public class CustomExceptionHandler(ILogger<CustomExceptionHandler> logger) : IExceptionHandler
{
    public async ValueTask<bool> TryHandleAsync(HttpContext httpContext, Exception exception, CancellationToken cancellationToken)
    {
        logger.LogError("Error Message: {Message}, Time of Occurrence: {Time}",
                        exception.Message, DateTime.UtcNow);

        var (detail, title, statusCode) = exception switch
        {
            InternalServerException => (exception.Message, exception.GetType().Name, StatusCodes.Status500InternalServerError),
            BadRequestException => (exception.Message, exception.GetType().Name, StatusCodes.Status400BadRequest),
            NotFoundException => (exception.Message, exception.GetType().Name, StatusCodes.Status404NotFound),
            ValidationException => (exception.Message, exception.GetType().Name, StatusCodes.Status422UnprocessableEntity),
            AlreadyExistException => (exception.Message, exception.GetType().Name, StatusCodes.Status409Conflict),
            ForbiddenException => (exception.Message, exception.GetType().Name, StatusCodes.Status403Forbidden),
            _ => (exception.Message, exception.GetType().Name, StatusCodes.Status500InternalServerError)
        };

        httpContext.Response.StatusCode = statusCode;

        var problemDetails = new ProblemDetails
        {
            Title = title,
            Status = statusCode,
            Detail = detail,
            Instance = httpContext.Request.Path,
        };

        problemDetails.Extensions.Add("traceId", httpContext.TraceIdentifier);

        if (exception is ValidationException validationException)
        {
            problemDetails.Extensions.Add("errors", validationException.Errors);
        }

        await httpContext.Response.WriteAsJsonAsync(problemDetails, cancellationToken: cancellationToken);
        return true;
    }
}
