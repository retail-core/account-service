using AccountService.API.Data.Entities;
using AccountService.API.Data.Repositories;

namespace AccountService.API.Services;

public interface IBusinessService
{
    Task<Business?> CreateBusinessAsync(string name, Guid ownerId);
}

public class BusinessService(IBusinessRepository _repo) : IBusinessService
{
    public async Task<Business?> CreateBusinessAsync(string name, Guid ownerId)
    {
        var business = new Business
        {
            Id = Guid.NewGuid(),
            Name = name,
            UserId = ownerId,
        };

        await _repo.CreateAsync(business);
        return business;
    }
}
