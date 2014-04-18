package resolver

type ResolverStorage interface {
    Get(string) string
    Set(string, string)
    Delete(string)
    List() []string
}
