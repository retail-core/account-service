using AccountService.API.Services;

namespace AccountService.API.Endpoints;

public static class StoreEndpoints
{
    public static void MapStoreEndpoints(this IEndpointRouteBuilder routes)
    {
        var group = routes.MapGroup("/v1/stores").WithTags("Store");

        group.MapGet("/{userId}", async (Guid userId, IStoreService _service) =>
        {
            var stores = await _service.GetStoresByUserIdAsync(userId);
            return Results.Ok(stores);
        });
    }
}
