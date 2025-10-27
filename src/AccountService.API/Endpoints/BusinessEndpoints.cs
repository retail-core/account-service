using AccountService.Core.Interfaces;
using AccountService.Data;

namespace AccountService.API.Endpoints;

public static class BusinessEndpoints
{
    public static void MapBusinessEndpoints(this IEndpointRouteBuilder routes)
    {
        var group = routes.MapGroup("/api/v1/business").WithTags("Business");

        group.MapGet("/", async (IBusinessService _service) =>
        {
            var businesses = await _service.GetAllAsync();
            return Results.Ok(businesses);
        });
    }
}
