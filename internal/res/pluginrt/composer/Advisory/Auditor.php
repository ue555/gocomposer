<?php declare(strict_types=1);
// gocomposer plugin runtime: the constants of Composer\Advisory\Auditor
// (Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Advisory/Auditor.php, MIT license, (c) Nils Adermann, Jordi
// Boggiano); audits are run by gocomposer.

namespace Composer\Advisory;

class Auditor
{
    public const FORMAT_TABLE = 'table';

    public const FORMAT_PLAIN = 'plain';

    public const FORMAT_JSON = 'json';

    public const FORMAT_SUMMARY = 'summary';

    public const FORMATS = [
        self::FORMAT_TABLE,
        self::FORMAT_PLAIN,
        self::FORMAT_JSON,
        self::FORMAT_SUMMARY,
    ];

    /** @deprecated Use ListPolicyConfig::AUDIT_IGNORE instead */
    public const ABANDONED_IGNORE = 'ignore';
    /** @deprecated Use ListPolicyConfig::AUDIT_REPORT instead */
    public const ABANDONED_REPORT = 'report';
    /** @deprecated Use ListPolicyConfig::AUDIT_FAIL instead */
    public const ABANDONED_FAIL = 'fail';

    /** Values to determine the audit result. */
    public const STATUS_OK = 0;
    public const STATUS_FAILED = 1;
}
