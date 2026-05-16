export interface AuthService {
  handler(request: Request): Response | Promise<Response>;
}
