<?php declare(strict_types=1);
// gocomposer plugin runtime: Composer\Command\BaseConfigCommand, standing in
// for Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/Command/BaseConfigCommand.php (MIT license, (c) Nils Adermann,
// Jordi Boggiano). Present so that plugins detecting the Composer version by
// this class (Composer 2.9+) see the version gocomposer implements; the
// config commands themselves are gocomposer's.

namespace Composer\Command;

abstract class BaseConfigCommand extends BaseCommand
{
}
