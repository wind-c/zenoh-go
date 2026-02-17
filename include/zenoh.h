#ifndef ZENOH_GO_ZENOH_H
#define ZENOH_GO_ZENOH_H

#ifdef _WIN32
typedef int z_result_t;
typedef void* z_config_t;
typedef void* z_session_t;
#define Z_OK 0

// Config functions
z_result_t z_config_new(z_config_t **config);
z_result_t z_config_default(z_config_t **config);
z_result_t zc_config_from_file(z_config_t *config, const char *path);
z_result_t zc_config_from_str(z_config_t *config, const char *s);
z_result_t zc_config_insert_json5(z_config_t config, const char *key, const char *value);
void z_config_drop(z_config_t config);

// Session functions
z_result_t z_open(z_session_t *session, z_config_t config, void *opts);
z_result_t z_close(z_session_t session, void *opts);
#else
#include_next <zenoh.h>
#endif

#endif
