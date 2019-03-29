
export interface IAuthPayload {
    ID: string;
    Roles: string[];
}

export interface IGlobalProps {
    authPayload: IAuthPayload;
    error?: string;
}