using AccountService.API.Data.Entities;
using AccountService.API.Data.Repositories;

namespace AccountService.API.Services;

public interface IStoreService
{
    Task<Store> CreateStoreAsync(string name, Guid businessId);
    Task<List<Store>> GetStoresByUserIdAsync(Guid userId);
}

public class StoreService(IStoreRepository _repo) : IStoreService
{
    public async Task<Store> CreateStoreAsync(string name, Guid businessId)
    {
        var store = new Store
        {
            Id = Guid.NewGuid(),
            Name = name,
            BusinessId = businessId,
        };

        await _repo.CreateStoreAsync(store);
        return store;
    }


    public async Task<List<Store>> GetStoresByUserIdAsync(Guid userId)
    {
        return await _repo.GetStoresByUserIdAsync(userId);
    }
}
