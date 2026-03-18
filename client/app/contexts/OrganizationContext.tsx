import React, { createContext, useContext } from "react";

interface OrganizationContextProps {
  id: number;
  name: string;
  slug: string;
  isOwner: boolean;
}

const OrganizationContext = createContext<OrganizationContextProps | undefined>(
  undefined,
);

export const OrganizationProvider: React.FC<
  OrganizationContextProps & { children: React.ReactNode }
> = ({ id, name, slug, children, isOwner }) => {
  return (
    <OrganizationContext.Provider value={{ id, name, slug, isOwner }}>
      {children}
    </OrganizationContext.Provider>
  );
};

export const useOrganizationContext = () => {
  const context = useContext(OrganizationContext);
  if (!context) {
    throw new Error(
      "useOrganizationContext must be used within a OrganizationProvider",
    );
  }
  return context;
};
