using AccountService.API.Data.Entities;
using Microsoft.EntityFrameworkCore;
using AccountService.API.Exceptions;

namespace AccountService.API.Data.Repositories;

public interface IStoreRepository
{
    Task CreateStoreAsync(Store store);
    Task<List<Store>> GetStoresByUserIdAsync(Guid userId);
}

public class StoreRepository(AppDbContext _dbContext) : IStoreRepository
{
    public async Task CreateStoreAsync(Store store)
    {
        await _dbContext.Stores.AddAsync(store);
        await _dbContext.SaveChangesAsync();
    }

    public async Task<List<Store>> GetStoresByUserIdAsync(Guid userId)
    {
        var business = await _dbContext.Businesses
            .AsNoTracking()
            .FirstOrDefaultAsync(b => b.UserId == userId)
            ?? throw new NotFoundException("No business found for the given user ID.");

        return await _dbContext.Stores
            .AsNoTracking()
            .Where(s => s.BusinessId == business.Id)
            .ToListAsync();
    }
}