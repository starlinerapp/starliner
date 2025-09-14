import { Configuration, RootApiFactory } from "~/api/client/generated";

const configuration: Configuration = {
  isJsonMime(mime: string): boolean {
    const jsonMime = new RegExp(
      // eslint-disable-next-line no-control-regex
      "^(application/json|[^;/ \t]+/[^;/ \t]+[+]json)[ \t]*(;.*)?$",
      "i",
    );
    return (
      mime !== null &&
      (jsonMime.test(mime) ||
        mime.toLowerCase() === "application/json-patch+json")
    );
  },
  basePath: "http://server:9090",
};

export const rootApiFactory = RootApiFactory(configuration);
