namespace AccountService.API.Exceptions;

public class AlreadyExistException : Exception
{
    public AlreadyExistException(string message) : base(message)
    {
    }
    
    public AlreadyExistException(string name, object key) 
            : base($"Entity \"{name}\" ({key}) already exist.")
    {
    }
}
