using Microsoft.EntityFrameworkCore;
using AccountService.API.Data.Entities;

namespace AccountService.API.Data;

public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options)
        : base(options) { }

    public DbSet<Business> Businesses => Set<Business>();
    public DbSet<Store> Stores => Set<Store>();
    public DbSet<Staff> Staffs => Set<Staff>();
    public DbSet<StaffStore> StaffStores => Set<StaffStore>();
    public DbSet<StaffPermission> StaffPermissions => Set<StaffPermission>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // Business → Stores
        modelBuilder.Entity<Business>()
            .HasMany(b => b.Stores)
            .WithOne(s => s.Business!)
            .HasForeignKey(s => s.BusinessId)
            .OnDelete(DeleteBehavior.Cascade);

        // Business → Staff
        modelBuilder.Entity<Business>()
            .HasMany(b => b.Staffs)
            .WithOne(s => s.Business!)
            .HasForeignKey(s => s.BusinessId)
            .OnDelete(DeleteBehavior.Cascade);

        // Staff ↔ Store (many-to-many)
        modelBuilder.Entity<StaffStore>()
            .HasKey(ss => new { ss.StaffId, ss.StoreId });

        modelBuilder.Entity<StaffStore>()
            .HasOne(ss => ss.Staff)
            .WithMany(s => s.StaffStores)
            .HasForeignKey(ss => ss.StaffId);

        modelBuilder.Entity<StaffStore>()
            .HasOne(ss => ss.Store)
            .WithMany(s => s.StaffStores)
            .HasForeignKey(ss => ss.StoreId);

        // Staff → Permissions
        modelBuilder.Entity<StaffPermission>()
            .HasOne(p => p.Staff)
            .WithMany(s => s.Permissions)
            .HasForeignKey(p => p.StaffId)
            .OnDelete(DeleteBehavior.Cascade);
    }
}
