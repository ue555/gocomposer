<?php declare(strict_types=1);
// gocomposer plugin runtime: copied unchanged from Composer f83588392df43d4fd027d84d26de459946ef4e80
// src/Composer/EventDispatcher/ScriptExecutionException.php (MIT license, (c) Nils Adermann, Jordi Boggiano).

/*
 * This file is part of Composer.
 *
 * (c) Nils Adermann <naderman@naderman.de>
 *     Jordi Boggiano <j.boggiano@seld.be>
 *
 * For the full copyright and license information, please view the LICENSE
 * file that was distributed with this source code.
 */

namespace Composer\EventDispatcher;

/**
 * Thrown when a script running an external process exits with a non-0 status code
 *
 * @author Jordi Boggiano <j.boggiano@seld.be>
 */
class ScriptExecutionException extends \RuntimeException
{
}
