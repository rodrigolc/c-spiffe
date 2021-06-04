#include "federation/src/watcher.h"
#include "spiffeid/src/trustdomain.h"
#include <check.h>

START_TEST(test_spiffebundle_Watcher_New)
{
    spiffebundle_Watcher *watcher = spiffebundle_Watcher_New();

    ck_assert_ptr_ne(watcher, NULL);
    ck_assert_ptr_ne(watcher->endpoints, NULL);

    spiffebundle_Watcher_Free(watcher);
}
END_TEST

START_TEST(test_spiffebundle_Watcher_AddHttpsWebEndpoint)
{
    spiffebundle_Watcher *watcher = spiffebundle_Watcher_New();
    spiffeid_TrustDomain td = { "example.com" };
    err_t err = spiffebundle_Watcher_AddHttpsWebEndpoint(
        watcher, "example.com/blabla", td);

    ck_assert_uint_eq(err, NO_ERROR);
    
    const int idx = shgeti(watcher->endpoints, td.name);
    
    ck_assert_int_ge(idx, 0);
    ck_assert_ptr_ne(watcher->endpoints[idx].value, NULL);

    spiffebundle_Watcher_Free(watcher);
}
END_TEST

START_TEST(test_piffebundle_Watcher_AddHttpsSpiffeEndpoint)
{
    spiffebundle_Watcher *watcher = spiffebundle_Watcher_New();
    spiffeid_TrustDomain td = { "example.com" };
    spiffebundle_Source *source = NULL;
    err_t err = piffebundle_Watcher_AddHttpsSpiffeEndpoint(
        watcher, "example.com/blabla", td, "spiffe://example.com/workload1",
        source);

    ck_assert_uint_eq(err, NO_ERROR);
    
    const int idx = shgeti(watcher->endpoints, td.name);
    
    ck_assert_int_ge(idx, 0);
    ck_assert_ptr_ne(watcher->endpoints[idx].value, NULL);

    spiffebundle_Watcher_Free(watcher);
    spiffebundle_Source_Free(source);
}
END_TEST

START_TEST(test_spiffebundle_Watcher_Start)
{
    spiffebundle_Watcher *watcher = spiffebundle_Watcher_New();
    err_t err = spiffebundle_Watcher_Start(watcher);

    ck_assert_uint_eq(err, NO_ERROR);

    spiffebundle_Watcher_Free(watcher);
}
END_TEST

START_TEST(test_spiffebundle_Watcher_Stop) 
{
    spiffebundle_Watcher *watcher = spiffebundle_Watcher_New();
    err_t err = spiffebundle_Watcher_Start(watcher);
    err = spiffebundle_Watcher_Stop(watcher);

    ck_assert_uint_eq(err, NO_ERROR);

    spiffebundle_Watcher_Free(watcher);
}
END_TEST

START_TEST(test_spiffebundle_Watcher_GetBundleForTrustDomain) {}
END_TEST

Suite *watcher_suite(void)
{
    Suite *s = suite_create("spiffebundle_watcher");
    TCase *tc_core = tcase_create("core");

    tcase_add_test(tc_core, test_spiffebundle_Watcher_New);
    tcase_add_test(tc_core, test_spiffebundle_Watcher_AddHttpsWebEndpoint);
    tcase_add_test(tc_core, test_piffebundle_Watcher_AddHttpsSpiffeEndpoint);
    tcase_add_test(tc_core, test_spiffebundle_Watcher_Start);
    tcase_add_test(tc_core, test_spiffebundle_Watcher_Stop);
    tcase_add_test(tc_core, test_spiffebundle_Watcher_GetBundleForTrustDomain);

    suite_add_tcase(s, tc_core);

    return s;
}

int main(int argc, char **argv)
{
    Suite *s = watcher_suite();
    SRunner *sr = srunner_create(s);

    srunner_run_all(sr, CK_NORMAL);
    const int number_failed = srunner_ntests_failed(sr);

    srunner_free(sr);

    return (number_failed == 0) ? EXIT_SUCCESS : EXIT_FAILURE;
}