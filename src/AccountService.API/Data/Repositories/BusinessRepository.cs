using AccountService.API.Data.Entities;
using AccountService.API.Data;
using Microsoft.EntityFrameworkCore;

namespace AccountService.API.Data.Repositories;

public interface IBusinessRepository
{
    Task CreateAsync(Business business);
    Task<Business?> GetBusinessByUserId(Guid userId);
}

public class BusinessRepository(AppDbContext _db) : IBusinessRepository  
{
    public async Task CreateAsync(Business business)
    {
        await _db.Businesses.AddAsync(business);
        await _db.SaveChangesAsync();
    }
    
    public async Task<Business?> GetBusinessByUserId(Guid userId)
    {
        var businesses = await _db.Businesses
                            .AsNoTracking()
                            .Where(b => b.UserId == userId)
                            .ToListAsync();
        return businesses.FirstOrDefault(b => b.UserId == userId);
    }
}