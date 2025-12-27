namespace AccountService.API.Endpoints;

public static class BusinessEndpoints
{
    public static void MapBusinessEndpoints(this IEndpointRouteBuilder routes)
    {
        var group = routes.MapGroup("/v1/businesses").WithTags("Business");
    }
}
