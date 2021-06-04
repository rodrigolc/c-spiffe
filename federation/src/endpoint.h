#ifndef INCLUDE_SPIFFEBUNDLE_ENDPOINT_H
#define INCLUDE_SPIFFEBUNDLE_ENDPOINT_H

#include "bundle/spiffebundle/src/bundle.h"
#include "bundle/spiffebundle/src/set.h"
#include "bundle/spiffebundle/src/source.h"
#include "spiffeid/src/id.h"
#include "spiffeid/src/trustdomain.h"
#include "utils/src/util.h"

#ifdef __cplusplus
extern "C" {
#endif

typedef struct spiffebundle_Source spiffebundle_Source;

typedef enum spiffebundle_Endpoint_Profile {
    NONE = 0,
    HTTPS_WEB,
    HTTPS_SPIFFE
} spiffebundle_Endpoint_Profile;

typedef struct spiffebundle_Endpoint {
    string_t url;
    profile;
    spiffeid_TrustDomain trust_domain;
    spiffeid_ID *spiffeID;
    spiffebundle_Bundle *bundle;
} spiffebundle_Endpoint;

spiffebundle_Endpoint *spiffebundle_Endpoint_New();
void spiffebundle_Endpoint_Free(spiffebundle_Endpoint *endpoint);

err_t spiffebundle_Endpoint_Config_HTTPS_WEB(
    spiffebundle_Endpoint *endpoint, string_t url,
    spiffeid_TrustDomain trust_domain);

err_t spiffebundle_Endpoint_Config_HTTPS_SPIFFE(
    spiffebundle_Endpoint *endpoint, string_t url,
    spiffeid_TrustDomain trust_domain, spiffeid_ID spiffeid,
    spiffebundle_Source *source);

err_t spiffebundle_Endpoint_Fetch(spiffebundle_Endpoint *endpoint);
spiffebundle_Bundle *spiffebundle_Endpoint_GetBundleForTrustDomain(
    spiffebundle_Endpoint *endpoint, spiffeid_TrustDomain trust_domain,
    err_t *err);
#ifdef __cplusplus
}
#endif

#endif // INCLUDE_SPIFFEBUNDLE_ENDPOINT_H
